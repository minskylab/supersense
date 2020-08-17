import React from "react";
import { ObserverBasicLayout } from "../layouts";

const ObserverPage: React.FC = () => {
    return (
        <div>
            <ObserverBasicLayout
                initialTitle="Your supersense Observer panel"
                initialMessage="#supersense"
            />
        </div>
    );
};

export { ObserverPage };
