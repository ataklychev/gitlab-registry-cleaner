package service

import (
	"sort"

	"github.com/ataklychev/gitlab-registry-cleaner/logger"
	"github.com/ataklychev/gitlab-registry-cleaner/repository"
	"github.com/xanzy/go-gitlab"
)

type GarbageCollectionServiceInterface interface {
	Run()
}

func NewGarbageCollectionService(
	threshold int,
	gitlabRepo repository.GitlabRepositoryInterface,
	logger logger.Logger,
) GarbageCollectionServiceInterface {
	// minimal number of saved images
	if threshold < 1 {
		threshold = 1
	}
	return &GarbageCollectionService{
		threshold:  threshold,
		gitlabRepo: gitlabRepo,
		logger:     logger,
	}
}

type GarbageCollectionService struct {
	threshold  int
	gitlabRepo repository.GitlabRepositoryInterface
	logger     logger.Logger
}

func (s *GarbageCollectionService) Run() {
	s.logger.Info("GarbageCollectionService.Run")

	projects := make(chan *gitlab.Project)

	go s.gitlabRepo.LoadProjects(projects)

	for project := range projects {
		s.cleanProject(project)
	}
}

func (s *GarbageCollectionService) cleanProject(project *gitlab.Project) {
	if nil == project {
		s.logger.Error("project required")

		return
	}

	s.logger.Infof("Clean project (%s)", project.PathWithNamespace)

	repos, err := s.gitlabRepo.ListRegistryRepositories(project.ID)
	if nil != err {
		s.logger.Error(err)

		return
	}

	s.cleanRepositories(project, repos)
}

func (s *GarbageCollectionService) cleanRepositories(
	project *gitlab.Project,
	repos []*gitlab.RegistryRepository,
) {
	for _, repo := range repos {
		s.cleanRepository(project, repo)
	}
}

func (s *GarbageCollectionService) cleanRepository(
	project *gitlab.Project,
	repo *gitlab.RegistryRepository,
) {
	tags, err := s.gitlabRepo.ListRegistryRepositoryTags(project.ID, repo.ID)
	if nil != err {
		s.logger.Error(err)

		return
	}

	if len(tags) > s.threshold {
		s.cleanTags(project, repo, tags)
	}
}

func (s *GarbageCollectionService) cleanTags(
	project *gitlab.Project,
	repo *gitlab.RegistryRepository,
	tags []*gitlab.RegistryRepositoryTag,
) {
	var err error

	// load details for each tag
	_tags := make([]*gitlab.RegistryRepositoryTag, len(tags))
	for i, tag := range tags {
		_tags[i], err = s.gitlabRepo.GetDetailsRegistryRepositoryTag(project.ID, repo.ID, tag.Name)
		if nil != err {
			s.logger.Error(err)

			return
		}
	}

	// sort tags by created_at DESC
	sort.Slice(_tags, func(i, j int) bool {
		return _tags[i].CreatedAt.After(*_tags[j].CreatedAt)
	})

	for i, tag := range _tags {
		if i >= s.threshold {
			s.logger.Infof("Delete %s", tag.Path)

			if err = s.gitlabRepo.DeleteRegistryRepositoryTag(project.ID, repo.ID, tag.Name); nil != err {
				s.logger.Error(err)
			}
		} else {
			s.logger.Infof("Safe %s", tag.Path)
		}
	}
}
