package service

type Service interface {
	repo *repository.Repository
	RecordUrlVisit(shortCode string) error
	GetUrlStats(shortCode string)
}

func (s *Service) RecordUrlVisit(shortCode string) error{
	visit := models.UrlVisit{
		shortCode: shortCode,
	}
	err := repo.SaveUrlVisit(visit)
	return err
}