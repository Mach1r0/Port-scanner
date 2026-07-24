# Port Scanner (Worker Pool) — Go

Scanner de portas TCP concorrente, construído para aprender e demonstrar o modelo de concorrência do Go (goroutines, channels, worker pool) aplicado a um caso real de reconhecimento de rede.

## Requisitos funcionais

- [ ] Receber um alvo (IP ou hostname) e um range de portas via flags de CLI
- [ ] Escanear portas TCP usando `net.DialTimeout`
- [ ] Implementar via **worker pool** (número de workers configurável), não goroutine-por-porta
- [ ] Agregar resultados via channel dedicado (`results channel`)
- [ ] Sincronizar finalização com `sync.WaitGroup`
- [ ] Output estruturado (tabela no terminal e/ou JSON) com: porta, status (open/closed/filtered), tempo de resposta
- [ ] Timeout configurável por conexão (evitar travar em portas com DROP silencioso)
- [ ] Modo verbose (mostra progresso) e modo silencioso (só resultado final)

## Requisitos não funcionais

- [ ] Rodar `go vet` e `golangci-lint` sem warnings
- [ ] Cobertura de testes nas funções principais (parsing de flags, worker, agregação de resultado)
- [ ] Detector de race conditions limpo (`go test -race`)
- [ ] Benchmark comparando execução sequencial vs. worker pool (documentado no README com números reais)

## Requisitos técnicos / flags sugeridas

```
-target       string  IP ou hostname alvo (obrigatório)
-ports        string  Range de portas, ex: "1-1000" ou "22,80,443" (default: "1-1024")
-workers      int     Número de workers concorrentes (default: 100)
-timeout      duration Timeout por conexão, ex: "500ms" (default: 1s)
-output       string  Formato de saída: "table" ou "json" (default: "table")
-verbose      bool    Mostra progresso em tempo real (default: false)
```

## Estrutura de pastas sugerida

```
port-scanner/
├── cmd/
│   └── scanner/
│       └── main.go
├── internal/
│   ├── scanner/
│   │   ├── worker.go
│   │   ├── scanner.go
│   │   └── scanner_test.go
│   └── output/
│       ├── table.go
│       └── json.go
├── benchmark/
│   └── results.md
├── go.mod
├── README.md
└── LICENSE
```

## Critérios de "pronto"

1. Escaneia 1000 portas em um host local em menos de alguns segundos com 100 workers
2. Testes passam com `-race` ativado
3. README final com benchmark comparativo (sequencial vs. worker pool) e print de execução
4. Código no GitHub com commits incrementais (não um único commit gigante) — isso importa para quem for revisar seu histórico

---

## README.md (rascunho para o repositório)

```markdown
# Port Scanner

Scanner de portas TCP concorrente escrito em Go, usando o padrão Worker Pool
para evitar exaustão de descritores de arquivo e maximizar throughput sem
sobrecarregar a rede alvo.

## Por que Worker Pool?

Escanear 65.535 portas sequencialmente é lento demais para uso prático.
Por outro lado, disparar uma goroutine por porta esgota os descritores de
arquivo do sistema operacional rapidamente. O padrão Worker Pool resolve
isso limitando o número de conexões concorrentes a um valor configurável,
mantendo alta performance sem instabilidade.

## Instalação

\`\`\`bash
git clone https://github.com/SEU_USUARIO/port-scanner
cd port-scanner
go build -o scanner ./cmd/scanner
\`\`\`

## Uso

\`\`\`bash
./scanner -target 192.168.1.1 -ports 1-1000 -workers 100
\`\`\`

## Benchmark

| Método       | Portas | Tempo   |
|--------------|--------|---------|
| Sequencial   | 1000   | Xs      |
| Worker Pool  | 1000   | Ys      |

(preencher com números reais após rodar)

## Arquitetura

[Job Channel] -> [N Workers] -> [Results Channel] -> [Agregador]

## Licença

MIT
```