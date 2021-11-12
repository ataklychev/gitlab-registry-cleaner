package repository

import (
	"log"
	"net/http"

	"github.com/ataklychev/gitlab-registry-cleaner/logger"
	"github.com/xanzy/go-gitlab"
)

func NewGitlabClient(accessToken string, baseAPIURL string) *gitlab.Client {
	git, err := gitlab.NewClient(accessToken, gitlab.WithBaseURL(baseAPIURL))
	if nil != err {
		log.Fatalf("Failed to create gitlab client: %v", err)
	}

	return git
}

type GitlabRepositoryInterface interface {
	// LoadProjects yield all projects
	LoadProjects(ch chan *gitlab.Project)
	// ListRegistryRepositories load container registry repositories
	ListRegistryRepositories(pid int) ([]*gitlab.RegistryRepository, error)
	// ListRegistryRepositoryTags load container registry repository tags
	ListRegistryRepositoryTags(pid, repository int) ([]*gitlab.RegistryRepositoryTag, error)
	// GetDetailsRegistryRepositoryTag load details of container registry repository tag
	GetDetailsRegistryRepositoryTag(pid, repository int, tagName string) (*gitlab.RegistryRepositoryTag, error)
	// DeleteRegistryRepositoryTag delete specified container registry repository tag
	DeleteRegistryRepositoryTag(pid, repository int, tagName string) error
}

func NewGitlabRepository(
	git *gitlab.Client,
	logger logger.Logger,
) GitlabRepositoryInterface {
	return &GitlabRepository{
		git:    git,
		logger: logger,
	}
}

type GitlabRepository struct {
	git    *gitlab.Client
	logger logger.Logger
}

// LoadProjects yield all projects to channel.
func (s *GitlabRepository) LoadProjects(ch chan *gitlab.Project) {
	maxPages := 1000
	perPage := 100
	page := 0

	for page < maxPages {
		options := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{PerPage: perPage, Page: page},
		}

		if projects, _, err := s.git.Projects.ListProjects(options); nil != err {
			s.logger.Error(err)
		} else {
			if len(projects) == 0 {
				close(ch)

				break
			} else {
				for _, project := range projects {
					ch <- project
				}
			}
		}

		page += 1
	}
}

func (s *GitlabRepository) ListRegistryRepositories(pid int) ([]*gitlab.RegistryRepository, error) {
	repos, resp, err := s.git.ContainerRegistry.ListRegistryRepositories(pid, &gitlab.ListRegistryRepositoriesOptions{})

	if resp.StatusCode == http.StatusForbidden {
		s.logger.Infof("Forbidden (%s)", pid)

		return repos, nil
	}

	return repos, err
}

func (s *GitlabRepository) ListRegistryRepositoryTags(pid, repository int) ([]*gitlab.RegistryRepositoryTag, error) {
	tags, _, err := s.git.ContainerRegistry.ListRegistryRepositoryTags(pid, repository, &gitlab.ListRegistryRepositoryTagsOptions{})

	return tags, err
}

func (s *GitlabRepository) GetDetailsRegistryRepositoryTag(pid, repository int, tagName string) (*gitlab.RegistryRepositoryTag, error) {
	tag, _, err := s.git.ContainerRegistry.GetRegistryRepositoryTagDetail(pid, repository, tagName)

	return tag, err
}

func (s *GitlabRepository) DeleteRegistryRepositoryTag(pid, repository int, tagName string) error {
	_, err := s.git.ContainerRegistry.DeleteRegistryRepositoryTag(pid, repository, tagName)

	return err
}
