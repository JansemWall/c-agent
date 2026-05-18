- UC01 - Registro de Nova Stack: Ao subir um docker-compose com a label ometris.enable=true, o agente detecta o novo container e envia o JSON de LLD (Low Level Discovery) para o Zabbix criar os itens automaticamente.

- UC02 - Monitoramento de Performance: A cada 30 segundos, o agente extrai as métricas de recursos do Docker e as envia para o Zabbix.

- UC03 - Alerta de Down-time: Se um container crítico (ex: DB) mudar para status unhealthy, o agente dispara um dado prioritário para o Zabbix gerar o alerta.