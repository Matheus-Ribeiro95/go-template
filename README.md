# Go Template

Código do template em Go.

    git clone https://github.com/Matheus-Ribeiro95/go-template.git

Tenha em mente que executando o código dessa forma, o arquivo irá chamar a função `panic`  por não estar conectado à database.

Para execução completa, antes execute a imagem da database:

    docker run -d -p 9042:9042 registry.gitlab.com/matheus-ribeiro95/public-images/template:scylladb

> Caso seja a primeira vez, aguarde um momento para a criação da imagem e importação dos dados padrões para a tabela `users` que o código usa

---

Instale as dependências:

    go mod download

Execute o código:

    go run main.go

Abra o link [localhost:8080/static](http://localhost:8080/static/)
