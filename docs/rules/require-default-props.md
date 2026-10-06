# Enforce a defaultProps definition for every prop that is not a required prop (react/require-default-props)

This rule aims to ensure that any non-required prop types of a class component
or `createReactClass` component has a corresponding `defaultProps` value.

> **Note**: You can provide types in runtime types using [PropTypes] and/or
statically using [TypeScript] or [Flow]. This rule will validate your prop types
regardless of how you define them.

`defaultProps` is deprecated for function components. Use ES6 default
parameters instead, and reserve `defaultProps` for class components and
`createReactClass` components.

## Rule Details

Examples of **incorrect** code for this rule:

```jsx
var Greeting = createReactClass({
  render: function() {
    return <div>Hello {this.props.foo} {this.props.bar}</div>;
  },

  propTypes: {
    foo: PropTypes.string,
    bar: PropTypes.string
  },

  getDefaultProps: function() {
    return {
      foo: "foo"
    };
  }
});
```

```jsx
class Greeting extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.foo} {this.props.bar}</h1>
    );
  }
}

Greeting.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};

Greeting.defaultProps = {
  foo: "foo"
};
```

Examples of **correct** code for this rule:

```jsx
function MyStatelessComponent({ foo, bar = 'some default' }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string.isRequired,
  bar: PropTypes.string
};
```

```jsx
class Greeting extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.foo} {this.props.bar}</h1>
    );
  }

  static propTypes = {
    foo: PropTypes.string,
    bar: PropTypes.string.isRequired
  };

  static defaultProps = {
    foo: "foo"
  };
}
```

```jsx
function MyStatelessComponent({ foo, bar }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string.isRequired,
  bar: PropTypes.string.isRequired
};
```

```jsx
function MyStatelessComponent({ foo = 'some default', bar = 'another default' }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};
```

```jsx
type Props = {
  foo: string,
  bar?: string
};

function MyStatelessComponent({ foo, bar = 'some default' }: Props) {
  return <div>Hello {foo} {bar}</div>;
}
```

```js
function NotAComponent({ foo, bar }) {}

NotAComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string.isRequired
};
```

## Rule Options

```js
...
"react/require-default-props": [<enabled>, { forbidDefaultForRequired: <boolean>, ignoreFunctionalComponents: <boolean> }]
...
```

* `enabled`: for enabling the rule. 0=off, 1=warn, 2=error. Defaults to 0.
* `forbidDefaultForRequired`: optional boolean to forbid prop default for a required prop. Defaults to false.
* `ignoreFunctionalComponents`: optional boolean to ignore this rule for functional components. Defaults to true. Set this to `false` to restore the deprecated behavior of requiring `defaultProps` on function
components.

### `forbidDefaultForRequired`

Forbids setting a default for props that are marked as `isRequired`.

Examples of **incorrect** code for this rule:

```jsx
class Greeting extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.foo} {this.props.bar}</h1>
    );
  }

  static propTypes = {
    foo: PropTypes.string,
    bar: PropTypes.string.isRequired
  };

  static defaultProps = {
    foo: "foo",
    bar: "bar"
  };
}
```

Examples of **correct** code for this rule:

```jsx
class Greeting extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.foo} {this.props.bar}</h1>
    );
  }

  static propTypes = {
    foo: PropTypes.string,
    bar: PropTypes.string.isRequired
  };

  static defaultProps = {
    foo: "foo"
  };
}
```

```jsx
function MyStatelessComponent({ foo, bar = 'bar' }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};
```

### `ignoreFunctionalComponents`

Functional components are ignored by default. Set this to `false` to enforce `defaultProps` for functional components as
well.

Examples of **incorrect** code for this rule:

```jsx
class Greeting extends React.Component {
  render() {
    return (
      <h1>Hello, {this.props.foo} {this.props.bar}</h1>
    );
  }

  static propTypes = {
    foo: PropTypes.string,
    bar: PropTypes.string.isRequired
  };

  static defaultProps = {
    foo: "foo",
    bar: "bar"
  };
}
```

Examples of **correct** code for this rule:

```jsx
function MyStatelessComponent({ foo, bar }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};
```

```jsx
const MyStatelessComponent = ({ foo, bar }) => {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};
```

```jsx
const MyStatelessComponent = function({ foo, bar }) {
  return <div>{foo}{bar}</div>;
}

MyStatelessComponent.propTypes = {
  foo: PropTypes.string,
  bar: PropTypes.string
};
```

## When Not To Use It

If you don't care about using `defaultProps` for your component's props that are not required, you can disable this rule.

# Resources
- [Official React documentation on defaultProps](https://facebook.github.io/react/docs/typechecking-with-proptypes.html#default-prop-values)

[PropTypes]: https://reactjs.org/docs/typechecking-with-proptypes.html
[TypeScript]: https://www.typescriptlang.org/
[Flow]: https://flow.org/
