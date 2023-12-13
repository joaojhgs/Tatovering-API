package models

type AgendamentoUsuario struct {
	Id               string `json:"id"`
	ClienteId        string `json:"cliente_id"`
	ServicoId        string `json:"servico_id"`
	TatuadorId       string `json:"tatuador_id"`
	EstudioId        string `json:"estudio_id"`
	Duracao          int    `json:"duracao"`
	Status           string `json:"status"`
	Email           string `json:"email"`
	Observacao       string `json:"observacao"`
	DataInicio       string `json:"data_inicio"`
	DataTermino      string `json:"data_termino"`
	DataCriacao      string `json:"data_criacao"`
	DataAtualizacao  string `json:"data_atualizacao"`
	DataCancelamento string `json:"data_cancelamento"`
}


type ViewUsuarioAgendamentosTatuador struct {
	Duracao          int    `json:"duracao"`
	DataInicio       string `json:"data_inicio"`
	DataTermino      string `json:"data_termino"`
}