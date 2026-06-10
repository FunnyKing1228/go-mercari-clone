# Mercari Clone Demo Client

This is a lightweight React/Vite frontend used to demonstrate the Go marketplace API.

It currently focuses on the item browsing experience:

- Fetches items from `http://localhost:8080/items`.
- Displays item name, price, status, and image when available.
- Uses TypeScript types for the item response shape.

The main project is backend-focused. This frontend is intentionally small and is included as a local demo client rather than a production-grade web application.

## Run Locally

Start the backend API first from the repository root:

```bash
docker compose up -d db redis
go run main.go
```

Then start the frontend:

```bash
cd frontend
npm install
npm run dev
```

The Vite dev server usually starts on `http://localhost:5173`.

## Scripts

```bash
npm run dev      # Start the Vite dev server
npm run build    # Type-check and build the frontend
npm run lint     # Run ESLint
npm run preview  # Preview the production build locally
```
