package models

type Estudio struct {
	Id                     string  `json:"id"`
	Nome                   string  `json:"nome"`
	Email                  string  `json:"email"`
	TaxaAgendamento        float64 `json:"taxa_agendamento"`
	Localizacao            string  `json:"localizacao"`
	Telefone               string  `json:"telefone"`
	Descricao              string  `json:"descricao"`
	Endereco               string  `json:"endereco"`
	ImagemPerfil           string  `json:"imagem_perfil"`
	ImagemCapa             string  `json:"imagem_capa"`
	ProprietarioId         string  `json:"proprietario_id"`
	HorarioDeFuncionamento *struct {
		Segunda []string `json:"segunda"`
		Terca   []string `json:"terca"`
		Quarta  []string `json:"quarta"`
		Quinta  []string `json:"quinta"`
		Sexta   []string `json:"sexta"`
		Sabado  []string `json:"sabado"`
		Domingo []string `json:"domingo"`
	} `json:"horario_funcionamento"`
	DiasFuncionamento *struct {
		Segunda bool `json:"segunda"`
		Terca   bool `json:"terca"`
		Quarta  bool `json:"quarta"`
		Quinta  bool `json:"quinta"`
		Sexta   bool `json:"sexta"`
		Sabado  bool `json:"sabado"`
		Domingo bool `json:"domingo"`
	} `json:"dias_funcionamento"`
}

type EstudioPost struct {
	Nome                   string  `json:"nome"`
	Email                  string  `json:"email"`
	TaxaAgendamento        float64 `json:"taxa_agendamento"`
	Localizacao            string  `json:"localizacao"`
	Telefone               string  `json:"telefone"`
	Descricao              string  `json:"descricao"`
	Endereco               string  `json:"endereco"`
	ImagemPerfil           string  `json:"imagem_perfil"`
	ImagemCapa             string  `json:"imagem_capa"`
	HorarioDeFuncionamento *struct {
		Segunda []string `json:"segunda"`
		Terca   []string `json:"terca"`
		Quarta  []string `json:"quarta"`
		Quinta  []string `json:"quinta"`
		Sexta   []string `json:"sexta"`
		Sabado  []string `json:"sabado"`
		Domingo []string `json:"domingo"`
	} `json:"horario_funcionamento"`
	DiasFuncionamento *struct {
		Segunda bool `json:"segunda"`
		Terca   bool `json:"terca"`
		Quarta  bool `json:"quarta"`
		Quinta  bool `json:"quinta"`
		Sexta   bool `json:"sexta"`
		Sabado  bool `json:"sabado"`
		Domingo bool `json:"domingo"`
	} `json:"dias_funcionamento"`
}

type Usuario struct {
	Id              string `json:"id"`
	Nome            string `json:"nome"`
	TelefoneCelular string `json:"telefone_celular"`
	Cpf             string `json:"cpf"`
	Rg              string `json:"rg"`
	Status          string `json:"status"`
	Endereco        string `json:"endereco"`
	Email           string `json:"email"`
}

type UsuarioView struct {
	Id              string `json:"id"`
	Nome            string `json:"nome"`
	TelefoneCelular string `json:"telefone_celular"`
	Cpf             string `json:"cpf"`
	Rg              string `json:"rg"`
	Status          string `json:"status"`
	Endereco        string `json:"endereco"`
	TatuadorId      string `json:"tatuador_id"`
	EstudioId       string `json:"estudio_id"`
	IsProprietario  bool   `json:"is_proprietario"`
}

type Tatuador struct {
	Id             string   `json:"id"`
	Nome           string   `json:"nome"`
	EstudioId      string   `json:"estudio_id"`
	Experiencia    int      `json:"experiencia"`
	EstiloTatuagem []string `json:"estilo_tatuagem"`
	Status         string   `json:"status"`
	Tipo           string   `json:"tipo"`
	ImgemPerfil    string   `json:"imagem_perfil"`
	ImgemCapa      string   `json:"imagem_capa"`
	RedesSociais   *struct {
		Instagram string `json:"instagram"`
		X         string `json:"x"`
		Facebook  string `json:"facebook"`
	} `json:"redes_sociais"`
}

type TatuadorPost struct {
	Nome           string   `json:"nome"`
	Experiencia    int      `json:"experiencia"`
	EstiloTatuagem []string `json:"estilo_tatuagem"`
	Status         string   `json:"status"`
	Tipo           string   `json:"tipo"`
	ImgemPerfil    string   `json:"imagem_perfil"`
	ImgemCapa      string   `json:"imagem_capa"`
	RedesSociais   *struct {
		Instagram string `json:"instagram"`
		X         string `json:"x"`
		Facebook  string `json:"facebook"`
	} `json:"redes_sociais"`
}

type Tatuagem struct {
	Id         string  `json:"id"`
	Imagem     string  `json:"imagem"`
	Preco      float64 `json:"preco"`
	Tamanho    int     `json:"tamanho"`
	Cor        string  `json:"cor"`
	Estilo     string  `json:"estilo"`
	TatuadorId string  `json:"tatuador_id"`
}

type TatuagensFavortas struct {
	Id         string  `json:"id"`
	Imagem     string  `json:"imagem"`
	Preco      float64 `json:"preco"`
	Tamanho    int     `json:"tamanho"`
	Cor        string  `json:"cor"`
	Estilo     string  `json:"estilo"`
	TatuadorId string  `json:"tatuador_id"`
}

type Favoritos struct {
	UsuarioId  string `json:"usuario_id"`
	TatuagemId string `json:"tatuagem_id"`
}
