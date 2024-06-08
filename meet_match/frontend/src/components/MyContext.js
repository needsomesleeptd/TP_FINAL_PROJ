import { createContext, useState } from "react";

const UserContext = createContext({});

const UserProvider = ({ children }) => {

    const [letsSwipe, setLetsSwipe] = useState(false);

    return (
        <UserContext.Provider value={{ letsSwipe, setLetsSwipe }}>
            {children}
        </UserContext.Provider>
    );
}

export {UserProvider, UserContext }
