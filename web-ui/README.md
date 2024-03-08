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
- [ ] Make back button in staking modal larger
- [ ] Fix skeleton spam when searching for non existent validator in staking modal

**Governance**

- [ ] Build liquid governance page

**UI/UX**

- [ ] Double check breakpoints

**Mobile Menu**

- [ ] improve mobile menu

**DevOps**

- [ ] Add doc for adding networks

**Blockers**

- [ ] Fix CCX in assets page
