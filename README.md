# Rinha de backend

Implementação em elixir da [rinha de backend 2023 Q3](https://github.com/zanfranceschi/rinha-de-backend-2023-q3).

Comecei a trabalhar no projeto depois de sairem todos os resultados e os vídeos do Akita e MrPowerGamerBR, então a maioria dos "truques" já eram amplamente conhecidos. O que fiz foi tentar explorar novos "truques" além dos que eram conhecidos, além de usar esse projeto como um tempo de estudos e exploração de novas tecnológias.

## Outras implementações

- [Elixir](https://github.com/ogabriel/rinha-de-backend-elixir)
- [Ruby](https://github.com/ogabriel/rinha-de-backend-ruby)
- [PHP](https://github.com/ogabriel/rinha-de-backend-php)

## Objetivo

Implementar a rinha em ruby on rails, e sem fazer nenuma adição de cache e batch insert.

Utilizando os aprendizados anteriores com a [implementação em elixir](https://github.com/ogabriel/rinha-backend-elixir) e [implementação em ruby](https://github.com/ogabriel/rinha-backend-ruby).

## Implementações testadas

1. brincar com as flags:
    a. GOGC
    b. GOMEMLIMIT
    c. GOMAXPROCS
2. passar uuid pro banco
3. tirar uuid pra apliação - _tira casaco, bota casado_
4. otimiza settings do postgresql
    - isso dadqui fez uma diferença fenomenal, principalmente desabilitar autovacuum e fsync

## Conclusões


## Resultados

### Desktop

|CPU|RAM|
|---|---|
|Ryzen 5900X|32GB|

#### Duas instâncias (com nginx)

##### Resultado do gatling navegador

![resultado gatling navegador part 1](./images/desktop/two/gatling-browser-1.png)
![resultado gatling navegador part 2](./images/desktop/two/gatling-browser-2.png)

##### Resultado do gatling console

```
Simulation RinhaBackendSimulation completed in 205 seconds
Parsing log file(s)...
Parsing log file(s) done
Generating reports...

================================================================================
---- Global Information --------------------------------------------------------
> request count                                     114991 (OK=114991 KO=0     )
> min response time                                      0 (OK=0      KO=-     )
> max response time                                     10 (OK=10     KO=-     )
> mean response time                                     0 (OK=0      KO=-     )
> std deviation                                          1 (OK=1      KO=-     )
> response time 50th percentile                          0 (OK=0      KO=-     )
> response time 75th percentile                          1 (OK=1      KO=-     )
> response time 95th percentile                          1 (OK=1      KO=-     )
> response time 99th percentile                          1 (OK=1      KO=-     )
> mean requests/sec                                558.209 (OK=558.209 KO=-     )
---- Response Time Distribution ------------------------------------------------
> t < 800 ms                                        114991 (100%)
> 800 ms <= t < 1200 ms                                  0 (  0%)
> t >= 1200 ms                                           0 (  0%)
> failed                                                 0 (  0%)
================================================================================
A contagem de pessoas é: 46576
```

#### Uma instância (sem nginx)

##### Resultado do gatling navegador

![resultado gatling navegador part 1](./images/desktop/one/gatling-browser-1.png)
![resultado gatling navegador part 2](./images/desktop/one/gatling-browser-2.png)

##### Resultado do gatling console

```
Simulation RinhaBackendSimulation completed in 205 seconds
Parsing log file(s)...
Parsing log file(s) done
Generating reports...

================================================================================
---- Global Information --------------------------------------------------------
> request count                                     114991 (OK=114991 KO=0     )
> min response time                                      0 (OK=0      KO=-     )
> max response time                                     10 (OK=10     KO=-     )
> mean response time                                     0 (OK=0      KO=-     )
> std deviation                                          1 (OK=1      KO=-     )
> response time 50th percentile                          0 (OK=0      KO=-     )
> response time 75th percentile                          1 (OK=1      KO=-     )
> response time 95th percentile                          1 (OK=1      KO=-     )
> response time 99th percentile                          1 (OK=1      KO=-     )
> mean requests/sec                                558.209 (OK=558.209 KO=-     )
---- Response Time Distribution ------------------------------------------------
> t < 800 ms                                        114991 (100%)
> 800 ms <= t < 1200 ms                                  0 (  0%)
> t >= 1200 ms                                           0 (  0%)
> failed                                                 0 (  0%)
================================================================================
A contagem de pessoas é: 46576
```

### Laptop

|CPU|RAM|
|---|---|
|Ryzen 4750U|16GB|

#### Duas instâncias (com nginx)

##### Resultado do gatling navegador

![resultado gatling navegador part 1](./images/laptop/two/gatling-browser-1.png)
![resultado gatling navegador part 2](./images/laptop/two/gatling-browser-2.png)

##### Resultado do gatling console

```
Simulation RinhaBackendSimulation completed in 205 seconds
Parsing log file(s)...
Parsing log file(s) done
Generating reports...

================================================================================
---- Global Information --------------------------------------------------------
> request count                                     114991 (OK=114991 KO=0     )
> min response time                                      0 (OK=0      KO=-     )
> max response time                                     41 (OK=41     KO=-     )
> mean response time                                     1 (OK=1      KO=-     )
> std deviation                                          1 (OK=1      KO=-     )
> response time 50th percentile                          1 (OK=1      KO=-     )
> response time 75th percentile                          1 (OK=1      KO=-     )
> response time 95th percentile                          2 (OK=2      KO=-     )
> response time 99th percentile                          2 (OK=2      KO=-     )
> mean requests/sec                                560.932 (OK=560.932 KO=-     )
---- Response Time Distribution ------------------------------------------------
> t < 800 ms                                        114991 (100%)
> 800 ms <= t < 1200 ms                                  0 (  0%)
> t >= 1200 ms                                           0 (  0%)
> failed                                                 0 (  0%)
================================================================================

Reports generated in 0s.
Please open the following file: file:///home/gabriel/workspace/rinha-de-backend-2023-q3/resultados/local/rinhabackendsimulation-20231108115309174/index.html
================================================================================
A contagem de pessoas é: 46576
```

##### Recusos do docker durante a parte mais pesada do teste

![Recusos do docker durante a parte mais pesada do teste](./images/laptop/two/docker-stats.png)

#### Uma instância (sem nginx)

##### Resultado do gatling navegador

![resultado gatling navegador part 1](./images/laptop/one/gatling-browser-1.png)
![resultado gatling navegador part 2](./images/laptop/one/gatling-browser-2.png)

##### Resultado do gatling console

```
```

##### Recusos do docker durante a parte mais pesada do teste

![Recusos do docker durante a parte mais pesada do teste](./images/laptop/one/docker-stats.png)
