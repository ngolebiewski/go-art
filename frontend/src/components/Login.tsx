import { type ChangeEvent, type FormEvent, useState } from "react";

export default function Login() {
  const [formData, setFormData] = useState({
    userName: "",
    email: "",
    password: "",
  });

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    alert(
      `Username: ${formData.userName}\nEmail: ${formData.email}\nPassword: ${formData.password}`
    );
  };

  const fields: { label: string; name: keyof typeof formData; type: string }[] = [
    { label: "Name", name: "userName", type: "text" },
    { label: "Email", name: "email", type: "email" },
    { label: "Password", name: "password", type: "password" },
  ];

  return (
    <>
    <h2>REGISTER / LOGIN</h2>
    <form onSubmit={handleSubmit}>
      {fields.map(({ label, name, type }) => (
        <label key={name}>
          {label}:
          <input
            type={type}
            name={name}
            value={formData[name]}
            onChange={handleChange}
          />
        </label>
      ))}
      <button type="submit">Submit</button>
    </form>
    </>
  );
}
