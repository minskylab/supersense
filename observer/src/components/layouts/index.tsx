import React, { useState, useEffect } from "react";
import { Box, Flex, Text } from "theme-ui";
import { useEventsStreamSubscription, Event } from "../../generated/graphql";
import EventCard from "../atoms/eventCard";
interface ObserverBasicLayoutProps {
    initialTitle: string;
    initialMessage: string;
}

const ObserverBasicLayout: React.FC<ObserverBasicLayoutProps> = ({
    initialTitle,
    initialMessage,
}: ObserverBasicLayoutProps) => {
    const [events, setEvents] = useState<Event[]>([]);

    const sub = useEventsStreamSubscription();

    useEffect(() => {
        console.log(sub.data);
        if (sub.data?.events) {
            const event = sub.data.events as Event;
            setEvents([event, ...events]);
        }
    }, [sub.data]);
    // comment test
    return (
        <Box>
            <Flex p={3} bg={"secondary"}>
                <Box sx={{ flex: 1 }}>
                    <Text sx={{ fontFamily: "heading" }}>SUPERSENSE</Text>
                </Box>
                <Flex sx={{ flex: [1, 2, 4] }}>
                    <Text
                        sx={{ fontFamily: "body", display: ["none", "block"] }}
                    >
                        {initialTitle}
                    </Text>
                </Flex>
                <Flex sx={{ flex: 1, justifyContent: "flex-end" }}>
                    <Text sx={{ fontFamily: "body" }}>{initialMessage}</Text>
                </Flex>
            </Flex>
            <Box
                p={3}
                style={{
                    display: "grid",
                    gridAutoFlow: "row dense",
                    gridTemplateColumns:
                        "repeat(auto-fit, minmax(250px, 350px))",
                    gridGap: "18px",
                    justifyItems: "center",
                    justifyContent: "center",
                }}
            >
                {events
                    // .sort((e1: Event, e2: Event) => {
                    //     return e2.emmitedAt - e1.emmitedAt;
                    // })
                    .map((event) => (
                        <EventCard key={event.id} event={event} />
                    ))}
            </Box>
        </Box>
    );
};

export { ObserverBasicLayout };
