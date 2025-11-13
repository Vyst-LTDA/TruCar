# Relatório de Migração de API: Python para Go

Este relatório detalha todos os endpoints e a lógica de negócios da API Python que ainda precisam ser migrados para a nova API em Go.

## Endpoints Pendentes de Migração

### 1. Admin (`admin.py`)

- **Rotas:**
    - `GET /organizations/`: Listar organizações (Super Admin).
    - `PUT /organizations/{org_id}`: Atualizar uma organização (Super Admin).
    - `GET /users/demo`: Listar usuários de demonstração (Super Admin).
    - `GET /users/all`: Listar todos os usuários (Super Admin).
    - `POST /users/{user_id}/activate`: Ativar a conta de um usuário (Super Admin).
    - `POST /users/{user_id}/impersonate`: Personificar um usuário (Super Admin).
- **Lógica de Negócios:**
    - Controle de acesso rigoroso para Super Admins.
    - Funções para ativar e personificar usuários, que exigem lógica de negócios específica.
- **Modelos e Schemas:**
    - `Organization`, `User`, `Token`.

### 2. Custos (`costs.py`)

- **Rotas:**
    - `GET /`: Listar todos os custos da organização com filtros de data.
- **Lógica de Negócios:**
    - Filtragem de custos por período.
- **Modelos e Schemas:**
    - `VehicleCost`.

### 3. Dashboard (`dashboard.py`)

- **Rotas:**
    - `GET /manager`: Obter dados para o dashboard do gestor.
    - `GET /driver`: Obter dados para o dashboard do motorista.
    - `GET /vehicles/positions`: Obter a geolocalização dos veículos em tempo real.
    - `GET /demo-stats`: Obter limites e uso da conta de demonstração.
- **Lógica de Negócios:**
    - Agregação de dados complexa para KPIs, classificações e alertas.
    - Lógica de painel separada para gestores e motoristas.
- **Modelos e Schemas:**
    - Múltiplos, incluindo `Report`, `User`, `Vehicle`, e schemas de dashboard personalizados.

### 4. Documentos (`documents.py`)

- **Rotas:**
    - `POST /`: Criar um novo documento com upload de arquivo.
    - `GET /`: Listar documentos com filtros.
    - `DELETE /{doc_id}`: Excluir um documento e o arquivo associado.
- **Lógica de Negócios:**
    - Manipulação de uploads de arquivos e armazenamento.
    - Exclusão de arquivos físicos do servidor.
- **Modelos e Schemas:**
    - `Document`.

### 5. Multas (`fines.py`)

- **Rotas:**
    - `POST /`: Criar uma nova multa e disparar uma notificação.
    - `GET /`: Listar multas (visão diferente para gestores e motoristas).
    - `PUT /{fine_id}`: Atualizar uma multa.
    - `DELETE /{fine_id}`: Excluir uma multa.
- **Lógica de Negócios:**
    - Criação de notificações em segundo plano.
    - Lógica de permissão para diferenciar a visibilidade dos dados.
- **Modelos e Schemas:**
    - `Fine`, `Notification`.

### 6. Ordens de Frete (`freight_orders.py`)

- **Rotas:**
    - Múltiplas rotas para gestores e motoristas, incluindo criação, atribuição, início e conclusão de paradas.
- **Lógica de Negócios:**
    - Máquina de estados complexa para o ciclo de vida das ordens de frete.
    - Lógica para motoristas se atribuírem a fretes e gerenciarem paradas.
- **Modelos e Schemas:**
    - `FreightOrder`, `StopPoint`, `Journey`.

### 7. GPS (`gps.py`)

- **Rotas:**
    - `POST /ping`: Receber um ping de localização de um dispositivo.
- **Lógica de Negócios:**
    - Atualização da localização do veículo em tempo real. Endpoint de alta frequência.
- **Modelos e Schemas:**
    - `Vehicle`, `LocationHistory`.

### 8. Implementos (`implements.py`)

- **Rotas:**
    - CRUD completo para implementos (reboques, etc.).
    - Rotas separadas para listar implementos disponíveis e todos os implementos.
- **Lógica de Negócios:**
    - Gerenciamento de um recurso separado que pode ser anexado a veículos.
- **Modelos e Schemas:**
    - `Implement`.

### 9. Leaderboard (`leaderboard.py`)

- **Rotas:**
    - `GET /`: Obter dados do placar de líderes.
- **Lógica de Negócios:**
    - Agregação de dados para classificar o desempenho dos motoristas.
- **Modelos e Schemas:**
    - `User`, e schemas de leaderboard personalizados.

### 10. Notificações (`notifications.py`)

- **Rotas:**
    - `GET /`: Listar notificações para o usuário logado.
    - `GET /unread-count`: Contar notificações não lidas.
    - `POST /{notification_id}/read`: Marcar uma notificação como lida.
    - `POST /trigger-alerts`: Disparar a verificação de alertas do sistema.
- **Lógica de Negócios:**
    - Sistema de notificação em tempo real para os usuários.
    - Lógica para verificações de alertas em todo o sistema.
- **Modelos e Schemas:**
    - `Notification`.

### 11. Peças (`parts.py`)

- **Rotas:**
    - CRUD completo para peças e itens de inventário, incluindo upload de imagens e faturas.
    - Rotas para gerenciar o status dos itens de inventário (disponível, em uso, descartado).
- **Lógica de Negócios:**
    - Gerenciamento de inventário complexo com histórico de transações.
    - Manipulação de uploads de arquivos.
- **Modelos e Schemas:**
    - `Part`, `InventoryItem`, `InventoryTransaction`.

### 12. Desempenho (`performance.py`)

- **Rotas:**
    - `GET /leaderboard`: Obter o ranking de desempenho dos motoristas.
- **Lógica de Negócios:**
    - Semelhante a `leaderboard.py`, mas com foco em métricas de desempenho.
- **Modelos e Schemas:**
    - `User`, e schemas de desempenho personalizados.

### 13. Gerador de Relatórios (`report_generator.py`)

- **Rotas:**
    - `POST /generate`: Gerar relatórios em PDF.
- **Lógica de Negócios:**
    - Geração de PDFs a partir de templates HTML.
    - Coleta de dados para preencher os relatórios.
- **Modelos e Schemas:**
    - `ReportRequest`, e vários outros dependendo do tipo de relatório.

### 14. Relatórios (`reports.py`)

- **Rotas:**
    - Vários endpoints para gerar diferentes tipos de relatórios (gerencial, desempenho do motorista, consolidado do veículo).
- **Lógica de Negócios:**
    - Agregação de dados complexa para diferentes tipos de relatórios.
- **Modelos e Schemas:**
    - Vários schemas de relatório personalizados.

### 15. Configurações (`settings.py`)

- **Rotas:**
    - `GET /fuel-integration`: Obter configurações de integração de combustível.
    - `PUT /fuel-integration`: Atualizar configurações de integração de combustível.
- **Lógica de Negócios:**
    - Gerenciamento seguro de chaves de API e segredos (criptografia).
- **Modelos e Schemas:**
    - `Organization`.

### 16. Telemetria (`telemetry.py`)

- **Rotas:**
    - `POST /report`: Receber e processar um pacote de dados de telemetria.
- **Lógica de Negócios:**
    - Processamento de dados de telemetria de veículos em tempo real.
- **Modelos e Schemas:**
    - `Vehicle`, `TelemetryPayload`.

### 17. Pneus (`tires.py`)

- **Rotas:**
    - `GET /vehicles/{vehicle_id}/tires`: Obter a configuração de pneus de um veículo.
    - `POST /vehicles/{vehicle_id}/tires`: Instalar um pneu em um veículo.
    - `PUT /tires/{tire_id}/remove`: Remover um pneu de um veículo.
    - `GET /vehicles/{vehicle_id}/removed-tires`: Obter o histórico de pneus removidos.
- **Lógica de Negócios:**
    - Gerenciamento do ciclo de vida dos pneus, incluindo instalação e remoção.
- **Modelos e Schemas:**
    - `Tire`, `Vehicle`.

### 18. Componentes do Veículo (`vehicle_components.py`)

- **Rotas:**
    - `GET /vehicles/{vehicle_id}/components`: Listar componentes instalados em um veículo.
    - `POST /vehicles/{vehicle_id}/components`: Instalar um novo componente.
    - `PUT /vehicle-components/{component_id}/discard`: Descartar um componente.
- **Lógica de Negócios:**
    - Gerenciamento do ciclo de vida de componentes do veículo.
- **Modelos e Schemas:**
    - `VehicleComponent`.

### 19. Custos do Veículo (`vehicle_costs.py`)

- **Rotas:**
    - `GET /`: Listar custos associados a um veículo.
    - `POST /`: Criar um novo registro de custo para um veículo.
- **Lógica de Negócios:**
    - Gerenciamento de custos específicos de um veículo.
- **Modelos e Schemas:**
    - `VehicleCost`.
