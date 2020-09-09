import React, { FC, ReactElement, Component } from "react";
import { Flex, Box, Text, useColorMode, useColorModeValue } from "@chakra-ui/core";
import { Settings, Moon, Sun } from "react-feather";

interface HeaderProps {
    brand?: Component | ReactElement | string;
    initialTitle: string;
    initialMessage: string;
    onSettings?: () => void;
}

const Header: FC<HeaderProps> = ({ initialTitle, initialMessage, onSettings, brand = "SUPERSENSE" }: HeaderProps) => {
    const { colorMode, toggleColorMode } = useColorMode();

    const bgColor = useColorModeValue("teal.200", "teal.800");

    const changeColorModeIcon = useColorModeValue(
        <Moon onClick={toggleColorMode} />,
        <Sun onClick={toggleColorMode} />
    );
    return (
        <Flex
            p={4}
            bg={bgColor}
            position={"fixed"}
            top={0}
            width={"100%"}
            zIndex={10}
            // borderBottomStyle={"solid"}
            // borderBottom={"1px"}
            // borderBottomColor={"black"}
        >
            <Box flex={1}>{brand}</Box>
            <Flex flex={[1, 2, 4]}>
                <Text display={["none", "block"]}>{initialTitle}</Text>
            </Flex>
            <Flex flex={1} justifyContent={"flex-end"}>
                <Text fontFamily={"body"}>{initialMessage}</Text>
            </Flex>
            <Flex ml={4}>
                <Settings onClick={onSettings} />
            </Flex>
            <Flex ml={4} _hover={{ cursor: "pointer" }}>
                {changeColorModeIcon}
            </Flex>
        </Flex>
    );
};

export default Header;
