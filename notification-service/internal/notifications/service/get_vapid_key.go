package service

func (s *NotificationService) GetVapidKey() string {
	return s.VapidPublicKey
}
