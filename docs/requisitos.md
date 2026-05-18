Funcionais

    Auto-Discovery: Identificar containers via Docker Socket que possuam labels específicas.

    Métricas Core: Coletar consumo de memória (bytes), CPU (%) e status de saúde (Healthcheck).

    Push-Based: Enviar dados via protocolo Zabbix Trapper (TCP 10051) para evitar abertura de portas de entrada nas VMs.

    Subpath Mapping: Associar o container ao seu subpath definido nas labels para facilitar a visão de "Negócio".

Não-Funcionais

    Baixo Footprint: Consumo máximo de 15MB de RAM por instância.

    Static Binary: Compilação única para facilitar o deploy via SCP ou Ansible.

    Resiliência: Cache local temporário em caso de perda de conexão com o Zabbix Server.