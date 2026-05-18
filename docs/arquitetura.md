Camada,Responsabilidade,Ferramenta/Lib
Provider (Docker),"Escaneia o socket, filtra containers por labels e extrai stats brutos.",docker/docker/client
Engine (Mapper),Converte os dados brutos em objetos JSON formatados para o Zabbix (LLD e Trapper).,encoding/json
Transporter (Zabbix),"Gerencia a conexão TCP, handshake do protocolo Zabbix e re-tentativas.",go-zabbix ou net (nativo)