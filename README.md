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

**Staking Page**

- figure out a better way to fit the custom weight Ui into the modal.

- need to handle number displays better, sometimes they show NaN or undefined but quickly render

- need to update any QS endpoints that will potentially be outdated with the coming of 1.4

**UI/UX**

- customize wallet connect modal

- threejs liquid metal sphere landing page

- Mobile breakpoints

**Defi**

- design

- build

**Assets**

- design

- build

**Mobile Menu**
