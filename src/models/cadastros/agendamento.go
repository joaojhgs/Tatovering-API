package models

type CadastroAgendamentoUsuario struct {
	Id         string `json:"id"`
	ClienteId  string `json:"cliente_id"`
	ServicoId  string `json:"servico_id"`
	TatuadorId string `json:"tatuador_id"`
	EstudioId  string `json:"estudio_id"`
	Duracao    int    `json:"duracao"`
	Status     string `json:"status"`
	Observacao string `json:"observacao"`
	DataInicio string `json:"data_inicio"`
}
