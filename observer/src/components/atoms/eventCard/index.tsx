import React, { useState, useEffect } from "react";
import { Box, Text, Flex, Image } from "theme-ui";
import type { Event } from "../../../generated/graphql";
import { GitHub, Twitter, Info, Share } from "react-feather";
import { format } from "timeago.js";

interface EventCardProps {
    event: Event;
}

const EventCard: React.FC<EventCardProps> = ({ event }: EventCardProps) => {
    const [timeAgo, setTimeAgo] = useState<string>(
        format(event.emmitedAt, "en_US"),
    );

    useEffect(() => {
        const timer = setInterval(() => {
            setTimeAgo(format(event.emmitedAt, "en_US"));
        }, 1000);

        return () => {
            window.clearInterval(timer);
        };
    }, []);

    return (
        <Box
            key={event.id}
            p={3}
            sx={{
                display: "flex",
                flexDirection: "column",
            }}
            style={{
                border: "solid 2px rgba(27,27,27,1)",
                borderRadius: "8px",
                WebkitBoxShadow: "5px 5px 0px 0px rgba(27,27,27,1)",
                MozBoxShadow: "5px 5px 0px 0px rgba(27,27,27,1)",
                boxShadow: "5px 5px 0px 0px rgba(27,27,27,1)",
            }}
        >
            {/* <Box>{event.id}</Box> */}
            <Flex sx={{ alignContent: "center" }}>
                <Box>
                    {event.sourceName === "twitter" && <Twitter size={18} />}
                    {event.sourceName === "github" && <GitHub size={18} />}
                    {event.sourceName === "dummy" && <Info size={18} />}
                </Box>
                {/* <Box ml={2}>{event.title}</Box> */}
                <Box ml={"auto"}>
                    <Text color={"black"} opacity={0.5} className="timeago">
                        {timeAgo}
                    </Text>
                </Box>
            </Flex>
            <Box py={4} sx={{ alignContent: "center" }}>
                <Text sx={{ textAlign: "center" }} as={"h3"}>
                    {event.message}
                </Text>
            </Box>
            <Flex sx={{ alignContent: "center" }} mt={"auto"}>
                <Image
                    src={event.actor.photo}
                    width={36}
                    height={36}
                    sx={{ borderRadius: "50%" }}
                />
                <Box
                    ml={3}
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                    }}
                >
                    <Text as={"p"} sx={{ alignContent: "center" }}>
                        {event.actor.username || event.actor.name}
                    </Text>
                </Box>
                <Box
                    ml={"auto"}
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                    }}
                >
                    <Share />
                </Box>
            </Flex>
        </Box>
    );
};

export default EventCard;
