import { useParams, useNavigate } from "react-router";
import { API_URL } from "./App";

export const ConfirmationPage = () => {
  let { token } = useParams();
  const redirect = useNavigate();

  const handleConfirmation = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT",
    });

    if (!response.ok) {
      console.error("Failed to confirm email");
    } else {
      redirect("/");
      console.log("Email confirmed successfully");
    }
  };

  return (
    <div>
      <h1>Confirmation Page</h1>
      <p>Please confirm your email address by clicking the button below.</p>
      <button onClick={handleConfirmation}>Confirm Email</button>
    </div>
  );
};
