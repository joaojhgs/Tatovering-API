package models

type Estudio struct {
	ProprietarioId       string  `json:"proprietario_uuid"`
	Nome                 string  `json:"nome"`
	Email                string  `json:"email"`
	HorarioDeFuncionamento *struct {
		Segunda []string `json:"segunda"`
		Terca   []string `json:"terca"`
		Quarta  []string `json:"quarta"`
		Quinta  []string `json:"quinta"`
		Sexta   []string `json:"sexta"`
	} `json:"horario_funcionamento"`
}

type Usuario struct {
	Nome            string `json:"nome"`
	TelefoneCelular string `json:"telefone_celular"`
	Cpf             string `json:"cpf"`
	Rg              string `json:"rg"`
	Status          string `json:"status"`
	Endereco        string `json:"endereco"`
}

type Tatuador struct {
	Experiencia    int    `json:"experiencia"`
	EstiloTatuagem string `json:"estilo_tatuagem"`
	Status         string `json:"status"`
	Tipo           string `json:"tipo"`
	ImgemPerfil           string `json:"imagem_perfil"`
	RedesSociais   *struct {
		Instagram string `json:"instagram"`
		X         string `json:"x"`
		Facebook  string `json:"facebook"`
	} `json:"redes_sociais"`
}

type Tatuagem struct {
	Desenho         string   `json:"desenho"`
	Preco           float64  `json:"preco"`
	Tamanho          int     `json:"tamanho"`
	Cor             string   `json:"cor"`
	Estilo          string   `json:"estilo"`
}