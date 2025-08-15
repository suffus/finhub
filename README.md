# FinHub CRM

A modern, AI-enabled sales and inbound marketing CRM system built with Go (backend) and Next.js (frontend).

## Features

- **User Management**: Multi-tenant user system with role-based access control
- **Contact Management**: Comprehensive contact and company management
- **Lead Management**: Lead tracking with status and temperature management
- **Sales Pipeline**: Deal management with customizable pipelines and stages
- **Marketing Tools**: Campaign management and lead source tracking
- **Content Management**: Web pages, blog posts, videos, and podcasts
- **Performance Tracking**: Sales goals and team performance metrics
- **AI Integration**: AI assistant for enhanced productivity
- **Dynamic Picklists**: Configurable dropdown selections with search and caching

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **ORM**: GORM with PostgreSQL
- **Authentication**: JWT with bcrypt password hashing
- **Database**: PostgreSQL 15+

### Frontend
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: TanStack Query (React Query)
- **Icons**: Lucide React
- **UI Components**: Custom components with shadcn/ui patterns

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

## Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd finhub
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PostgreSQL: localhost:5432

### Local Development

1. **Backend Setup**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

2. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

3. **Database Setup**
   ```bash
   # Start PostgreSQL
   docker run -d \
     --name finhub-postgres \
     -e POSTGRES_DB=finhub \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 \
     postgres:15-alpine
   ```

## Environment Variables

### Backend
```bash
DATABASE_URL=postgres://postgres:password@localhost:5432/finhub?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
DEBUG=true
PORT=8080
```

### Frontend
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update current user

### Companies
- `GET /api/companies` - List companies
- `POST /api/companies` - Create company
- `GET /api/companies/:id` - Get company
- `PUT /api/companies/:id` - Update company
- `DELETE /api/companies/:id` - Delete company

### Contacts
- `GET /api/contacts` - List contacts
- `POST /api/contacts` - Create contact
- `GET /api/contacts/:id` - Get contact
- `PUT /api/contacts/:id` - Update contact
- `DELETE /api/contacts/:id` - Delete contact

### Leads
- `GET /api/leads` - List leads
- `POST /api/leads` - Create lead
- `GET /api/leads/:id` - Get lead
- `PUT /api/leads/:id` - Update lead
- `DELETE /api/leads/:id` - Delete lead

### Deals
- `GET /api/deals` - List deals
- `POST /api/deals` - Create deal
- `GET /api/deals/:id` - Get deal
- `PUT /api/deals/:id` - Update deal
- `DELETE /api/deals/:id` - Delete deal

### Picklists
- `GET /api/picklists/:entity` - Get picklist items (industries, companysizes, leadstatuses, leadtemperatures)
- `POST /api/picklists/search` - Search picklist items with pagination

## Database Schema

The system uses a comprehensive database schema with the following main entities:

- **Users & Tenants**: Multi-tenant user management
- **Companies & Contacts**: Business relationship management
- **Leads**: Lead tracking and qualification
- **Deals**: Sales pipeline management
- **Marketing**: Campaign and asset management
- **Tasks & Communications**: Activity tracking
- **Custom Fields**: Extensible data model

## Development

### Backend Development

1. **Install Go dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

2. **Seed the database with sample data**
   ```bash
   make setup
   # Or run individually:
   make build-seed
   make seed
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Run with hot reload**
   ```bash
   go install github.com/cosmtrek/air@latest
   air
   ```

### Picklist System

The application includes a comprehensive picklist system for dynamic dropdown selections:

- **Backend**: RESTful API endpoints with search and pagination
- **Frontend**: Reusable PicklistSelect component with search capabilities
- **Caching**: In-memory caching for performance optimization
- **Multi-select**: Support for both single and multiple selection modes

For detailed information, see [PICKLIST_README.md](PICKLIST_README.md).

### Frontend Development

1. **Install dependencies**
   ```bash
   cd frontend
   npm install
   ```

2. **Run development server**
   ```bash
   npm run dev
   ```

3. **Build for production**
   ```bash
   npm run build
   npm start
   ```

## Deployment

### AWS Deployment

1. **Build and push Docker images**
   ```bash
   docker build -t finhub-backend ./backend
   docker build -t finhub-frontend ./frontend
   ```

2. **Deploy using AWS ECS or EKS**
   - Use the provided Docker Compose file as a reference
   - Configure environment variables for production
   - Set up proper SSL certificates and load balancers

### Environment Configuration

For production deployment, ensure you have:

- Strong JWT secret
- Production database with proper security
- HTTPS enabled
- Proper CORS configuration
- Monitoring and logging setup

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please open an issue in the GitHub repository. 