import React from "react";
import { NextPage } from "next";
import { ObserverBasicLayout } from "@app/components/layouts";

const IndexPage: NextPage = () => {
    const eventsBuffer = 5 * 3;

    return (
        <div>
            <ObserverBasicLayout
                initialTitle="Your supersense Observer panel"
                initialMessage="#supersense"
                bufferSize={eventsBuffer}
            />
        </div>
    );
};

export default IndexPage;
