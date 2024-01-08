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

### Development ToDo

**Governance**

- add liquid staked governance (when its built)

**UI/UX**

- focus on mobile landscape breakpoints, (mainly staking page)

**Mobile Menu**

- connect wallet button

- graphic elements

- font size / style / decorations

**DevOps**

- make onboarding networks seamless

**Assets Page**

- claim rewards claim.test.quicksilver.zone/address/current \*/epoch

- intent query

- fix ibc transactions and ibc transaction amino registry

**Staking Page**

- check memo intent creation

- validator route app.quicksilver.zone/staking/chainId/valoperAddress

**Defi Page**

- fix some of the links
