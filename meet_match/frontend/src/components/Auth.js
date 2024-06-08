const RegisterRequest = async (name, age, gender, login, password) => {
    try {
        const response = await fetch('/api/user/SignUp', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              'user': {
                'name': name,
                'age': Number(age),
                'gender': gender,
                'login': login,
                'password': password
              }
            })
        });
  
        if (!response.ok) {
          const errorMessage = await response.text();
          throw new Error(errorMessage);
        }

        return await response.json();
  
      } catch (error) {
        console(error.message);
      }
}

const LoginRequest = async (login, password) => {
    try {
        const response = await fetch('/api/user/SignIn', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'login': login,
                'password': password
            })
        });

        if (!response.ok) {
          const errorMessage = await response.text();
          throw new Error(errorMessage);
        }

        return await response.json();

    } catch (error) {
        console(error.message);
    }
}

export { LoginRequest, RegisterRequest };
