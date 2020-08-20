import React, { useState, useEffect } from "react";
import { Box, Flex, Text } from "theme-ui";
import { useEventsStreamSubscription, Event } from "../../generated/graphql";
import EventCard from "../atoms/eventCard";



interface ObserverBasicLayoutProps {
    initialTitle: string;
    initialMessage: string;
    eventsBuffer?: number;
}

const ObserverBasicLayout: React.FC<ObserverBasicLayoutProps> = ({
    initialTitle,
    initialMessage,
    eventsBuffer=45,
}: ObserverBasicLayoutProps) => {
    const [events, setEvents] = useState<Event[]>([]);

    const sub = useEventsStreamSubscription();

    useEffect(() => {
        console.log(sub.data);
        if (sub.data?.eventStream) {

            const event = sub.data.eventStream as Event;
            if (events.length + 1 > eventsBuffer) {
                events.pop()
            }
            setEvents([event, ...events]);
        }
    }, [sub.data]);


    return (
        <Box>
            <Flex p={3} bg={"primary"}>
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
                    .map((event: Event) => (
                        // <animated.div key={key} style={{ ...props, width: "100%", height: "100%" }}>
                            <EventCard key={event.id} event={event} />
                        // </animated.div>
                    ))}
            </Box>
        </Box>
    );
};

export { ObserverBasicLayout };
