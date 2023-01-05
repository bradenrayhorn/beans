import { useQueries } from "@/constants/queries";
import { routes } from "@/constants/routes";
import { useOnLogout } from "@/context/AuthContext";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

export const useLogout = () => {
  const queries = useQueries({});
  const navigate = useNavigate();
  const onLogout = useOnLogout();

  const mutation = useMutation(queries.logout, {
    onSuccess: () => {
      onLogout();
      navigate(routes.login);
    },
  });

  return { ...mutation };
};
