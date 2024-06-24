# Multitheading CEP

Uma aplicação que usa multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:
- https://brasilapi.com.br/api/cep/v1/
- http://viacep.com.br/ws/


## Estrutura

- **AddressViaCep** e **AddressBrasilApi**: As duas estruturas separadas para representar os dados retornados pelas APIs http://viacep.com.br e https://brasilapi.com.br, respectivamente. Cada estrutura possui campos correspondentes aos dados retornados por cada API.

- **Channels e Tipos de Resultado**: Um canal chan interface{} para permitir que as goroutines enviem resultados de tipos diferentes ( *AddressViaCep ou *AddressBrasilApi) para a função principal.

- **Funções de Busca (fetchFromAPI1 e fetchFromAPI2)**: Essas funções são utilizadas para chamar as funções fetchAddressViaCep e fetchAddressBrasilApi, respectivamente, e enviar o resultado para o canal apropriado.

- **Exibição do Resultado**: A função displayAddress foi adicionada para formatar e exibir os resultados obtidos, indicando qual API forneceu os dados.

## Execução do Programa

Para executar o programa, use o comando no terminal:

```
go run main.go -cep {{CEP_DESEJADO}}
```