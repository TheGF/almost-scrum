

import theme, { Theme } from "@chakra-ui/theme";
import { mode } from "@chakra-ui/theme-tools";

const backgroundLight = "blue.100"
const backgroundDark = "gray.800"
const colorLight = "blue.900"
const colorDark = "gray.500"

const styles = {
    ...theme.styles,
    global: (props) => ({
        ...theme.styles.global,
        fontFamily: "body",
        fontWeight: "light",
        "html, body": {
            background: props.colorMode === "dark" ? backgroundDark : backgroundLight,
            color: props.colorMode === "dark" ?  colorDark : colorLight,
            lineHeight: "tall",
        },
        ".react-tags__search-input": {
            backgroundColor: props.colorMode === "dark" ? backgroundDark : backgroundLight,
        },
        ".mde-text": {
//            minHeight: "12em",
            backgroundColor: props.colorMode === "dark" ? "gray.700" : "white",
            color: props.colorMode === "dark" ? "yellow.200" : "black",
        },
        ".mde-preview-content": {
//            minHeight: "14em",
            backgroundColor: props.colorMode === "dark" ? "gray.700" : "white",
            color: props.colorMode === "dark" ? "yellow.200" : "black",
        },
        ".mde-header": {
            backgroundColor: props.colorMode === "dark" ? "gray.400" : "gray.100",
        },
        ".image-input": {
            backgroundColor: props.colorMode === "dark" ? "gray.400" : "gray.100",
        },
        ".task-viewer": {
//            maxHeight: "20em",
            overflow: "auto",
        }
    }),
};

const customTheme = {
    ...theme,
    fonts: {
        ...theme.fonts,
        body: `"Source Sans Pro",${theme.fonts.body}`,
        heading: `"Source Sans Pro",${theme.fonts.heading}`,
    },
    styles,
};

export default customTheme;
