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

    // const transitions = useTransition<Event, React.CSSProperties>(events, event=>event.id, {
    //     from: { transform: 'translate3d(0,-40px,0)' },
    //     enter: { transform: 'translate3d(0,0px,0)' },
    //     leave: { transform: 'translate3d(0,-40px,0)' },
    // });

    const sub = useEventsStreamSubscription();

    useEffect(() => {
        console.log(sub.data);
        if (sub.data?.events) {
            const event = sub.data.events as Event;
            setEvents([event, ...events]);
        }
    }, [sub.data]);


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
