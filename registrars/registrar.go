package registrars

type Registrar interface {
    CheckDomainAvailability(domain string) (bool, float64, error)
}