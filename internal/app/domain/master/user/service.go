package user

// UserService is a Domain Service.
// Digunakan ketika business rule melibatkan lebih dari satu entity,
// atau membutuhkan akses ke Repository yang tidak natural dimiliki oleh Entity itu sendiri.
//
// Contoh: Entity User tidak tahu apakah email-nya unik di sistem —
// dia hanya tahu data dirinya sendiri.
// Maka pengecekan uniqueness email adalah tanggung jawab Domain Service.
type UserService struct {
	Repo Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{Repo: repo}
}

// CheckEmailAvailability memastikan email belum dipakai user lain.
// Ini adalah domain rule, bukan application/orchestration logic.
func (s *UserService) CheckEmailAvailability(email string) error {
	existing, _ := s.Repo.FindByEmail(email)
	if existing != nil {
		return ErrEmailAlreadyRegistered
	}
	return nil
}
