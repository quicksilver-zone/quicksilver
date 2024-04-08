This is the Quicksilver Liquid Staking App bootstrapped with [`create-cosmos-app`](https://github.com/cosmology-tech/create-cosmos-app) and re-templated to work with Bun.

## Getting Started

First, install the packages and run the development server:

```bash
bun install && bun run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the webpage.

## Making Contributions

Please work on a branch with a title that reflects what you aim to contribute and open a pull request to the `main` branch.

### Dependencies

Please use this Prettier config to format your code before opening a pull request.

```
{
  "semi": true,
  "trailingComma": "all",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2
}

```

Please ensure your IDE is configured to use Typescript v4.9.3

### TODO:

**Staking**

- [ ] Add support for Dymension
- [ ] Add more weight options IE `equal`, `custom`, `most votes`, `lowest commission` etc

**Governance**

**UI/UX**

- [ ] refetch on tx success

**Mobile Menu**

- [ ] improve mobile menu

**DevOps**

- [ ] Finish doc for adding networks

**Has Blockers**

- [ ] Build liquid governance page

**Assets**

- [ ] Fix the way queries for networks and entries in components are created. Rather than defining one for each network, create a function that iterates through the .env entry or liveNetworks call for live networks and creates the queries and components for each. `pages/assets.tsx`
