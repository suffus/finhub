We are aiming to build a sales and inbound marketing focused CRM to compete with salesforce and hubspot that is AI enabled and is to be at least 90% written by AI.  

OK I want you to build a very simple initial project with a React front end and a Go backend with postgres as the database.

The tech stack for the frontend should be Next.js, React, TanStack Query, Tailwind + shadcn/ui.  The programming language shall be typescript.

The backend stack is Go with gorm and gin.  The database will be postgres and we shall deploy on AWS via docker.  For now use a simple username and password authentication and use JWT bearer tokens.

The initial UI design is attached in an image and specified further in frontend.md.

The initial data model which was built for Prisma is attached in the file schema.prisma.  We want to use Gorm, so you will have to translate this to Go!

The frontend of the sales and marketing management system has a common top toolbar with widgeths for search, configuration, alerts, an AI assistant, and the user's profile and account settings.

Along the left hand side is a toolbar with an extensible set of bookmarks, followed by menus for all the main modules - Candidates/Contacts/Leads/Deals, the Marketing System (campaign management, advertising, website lead reporting), Marketing content artifacts including web pages, blog articles, podcasts, videos, and the like, and a Sales management system used to manage teams of sales executives, set goals for them, measure progress and performance.

When viewing an object (such as Contacts or Companies or Leads, etc) there may be multiple views of the data (which are essentially filters over the database tables combined with visualizations of the results), and each individual object can be viewed in panel using templated views.

Please set up initial projects for the frontend and the backend (in the separate directories provided) and build a simple frontend that looks like the attached image in frontend.avif with a search bar and tool bar on the top, and a menu on the left hand side that will enable access to the various modules of the CRM.  Please set the database we have provided up in GORM and provide basic CRUD functions for the database objects we have provided, plus implement a simple user registration and login API call to support the creation of sessions.



