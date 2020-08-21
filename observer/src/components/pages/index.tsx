import React from "react";
import { ObserverBasicLayout } from "@app/components/layouts";

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
