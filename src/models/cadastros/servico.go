package models

type CadastroServicoUsuario struct {
	Id                string  `json:"id"`
	ClienteId         string  `json:"cliente_id"`
	TatuadorId        string  `json:"tatuador_id"`
	EstudioId         string  `json:"estudio_id"`
	TatuagemId        string  `json:"tatuagem_id"`
	Tipo              string  `json:"tipo"`
	Objetivo          string  `json:"objetivo"`
	Descricao         string  `json:"descricao"`
	Valor             float64 `json:"valor"`
	QtdeSessoes       int64   `json:"qtde_sessoes"`
	ImagemReferencia  string  `json:"imagem_referencia"`
	ImagemReferencia2 string  `json:"imagem_referencia2"`
	ImagemReferencia3 string  `json:"imagem_referencia3"`
}
