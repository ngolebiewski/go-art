## React + TypeScript SPA Cheatsheet

# Events
```tsx // Input change
const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
console.log(e.target.value);
};

// Form submit
const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
e.preventDefault();
};

// Button click
const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
console.log("clicked");
};
// Keyboard
const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
if (e.key === "Enter") console.log("Enter pressed");
};
```

# State
```tsx const [count, setCount] = useState<number>(0);
type User = { id: number; name: string; email: string };
const [user, setUser] = useState<User | null>(null);
const [items, setItems] = useState<string[]>([]);
```

# Props
```tsx 
type ButtonProps = {
label: string;
onClick: () => void;
function Button({ label, onClick }: ButtonProps): JSX.Element {
return <button onClick={onClick}>{label}</button>;
};
}
```

# Components

```tsx 
const Header = (): JSX.Element => <h1>My App</h1>;
```

# Context + Reducer
```tsx 
type State = { count: number };
type Action = { type: "inc" } | { type: "dec" };
function reducer(state: State, action: Action): State {
switch (action.type) {
case "inc": return { count: state.count + 1 };
case "dec": return { count: state.count - 1 };
}
}
const [state, dispatch] = useReducer(reducer, { count: 0 });
```

# Refs
```tsx 
const inputRef = useRef<HTMLInputElement | null>(null);
// Safe usage
if (inputRef.current) inputRef.current.focus();
```

# API Data
```tsx 
type User = { id: number; name: string; email: string };
async function fetchUser(): Promise<User> {
const res = await fetch("/api/user");
return res.json();
}
```

# Utility Types
```tsx 
type User = { id: number; name: string; email: string };
type UserPreview = Pick<User, "id" | "name">;
type UserWithoutEmail = Omit<User, "email">;
type PartialUser = Partial<User>; // all optional
type UserRecord = Record<string, User>;
```

# Common Gotchas
```tsx 
- "Object is possibly 'null'" → check before using.
if (user) console.log(user.name);
- "Type 'undefined' is not assignable" → initialize state or allow undefined.
- Generics with hooks → useState<Type>().
```

# Quick Checklist
[ ] Type events (ChangeEvent, FormEvent, etc.)
[ ] Type state explicitly (useState<Type>)
[ ] Create prop interfaces (Props)
[ ] Return JSX.Element from components
[ ] Type context + reducer (State, Action)
[ ] Use typed refs (useRef<HTMLInputElement>)
[ ] Define API response types (type User = {})
[ ] Use utility types (Partial, Pick, Omit, etc.)
[ ] Handle null/undefined safely