const ashOptions = localStorage.getItem("ash-options")
const uiEffects = !ashOptions || ashOptions.includes('ui-effect')

const colors = { 
    dark: {
        background: 'linear-gradient(#404e50, #101010)',   
        color: 'gray.800',    
        input: 'gray.200',
        panel1bg: 'linear-gradient(#88ceca, #689e9a)',
        panel2bg: 'linear-gradient(#a1a1a1, #919191)'
    },
    light: {
        background: 'linear-gradient(to right bottom, #a0aeb0, #606070)',
        color: 'gray.900',    
        input: 'white',
        panel1bg: 'linear-gradient(#e8ceca, #b89e9a)',
        panel2bg: 'linear-gradient(#f1f1f1, #c1c1c1)',
    }
}


const colorsFlat = { 
    dark: {
        background: 'gray.600',   
        color: 'gray.100',    
        input: 'gray.600',
        panel1bg: '#e8ceca',
        panel2bg: '#e1e1e1'
    },
    light: {
        background: '#909ea0',
        color: 'gray.900',    
        input: 'white',
        panel1bg: '#e8ceca',
        panel2bg: '#f1f1f1'
    }
}


export default function(mode, color) {
    return uiEffects ? colors[mode][color] : colorsFlat[mode][color]
} 