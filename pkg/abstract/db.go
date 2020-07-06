package abstract

type Line struct {
	RemoteHost string
	Time       time.Time
	Request    string
	Status     int
	Bytes      int
	Referer    string
	UserAgent  string
	URL        string
}

type DBLineInserter interface {
	Insert(ctx context.context, line Line) error
}

type DB interface {
	DBLineInserter
	Close()
}
