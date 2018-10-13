package domain

type Config struct {
	GrpcPort    int    `config:"default=9000;usage="`
	UseTLS      bool   `config:"default=false;usage="`
	SSLCertFile string `config:"default=ssl.crt;usage=SSL public key file"`
	SSLKeyFile  string `config:"default=ssl.key;usage=SSL private key file"`
}
