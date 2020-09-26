package utils

import "log"

//Check checa erros e responde de acordo com a criticidade
// crit - 0 erro esperado em alguns casos de uso, deve ser logado e retornado para tratamento.
// crit - 1 erro inesperado ou falha, nao deve entrar em produção, interrompe a execução e deve ser logado.
func Check(err error, crit int) error {
	//caso de uso com maioria de instâncias deve ser tratado primeiro para evitar execução de código desnecessária
	if err == nil {
		return nil
	}
	if crit == 0 {
		log.Println("Erro: ", err)

	}
	if crit == 1 {
		log.Fatal("Erro: ", err)
		panic(0)

	}
	return err
}
