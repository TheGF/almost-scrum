const ashOptions = localStorage.getItem("ash-options")
const uiEffects = !ashOptions || ashOptions.includes('ui-effect')

const colors = { 
    dark: {
        background: 'linear-gradient(to right bottom, #404e50, #101010)',   
        color: 'gray.800',    
        input: 'gray.200',
        panel1bg: 'linear-gradient(to right bottom, #88ceca, #689e9a)',
        panel2bg: 'linear-gradient(to right bottom, #a1a1a1, #919191)'
    },
    light: {
        background: 'linear-gradient(to right bottom, #909ea0, #606070, #303040)',
        color: 'gray.900',    
        panel1bg: 'linear-gradient(to right bottom, #3182ce, #2172be)',
//        panel1bg: 'linear-gradient(to right bottom, #51a2fe, #4192ee)',
        panel1color: 'white',
        yellowPanelBg: 'linear-gradient(to right bottom, #ecc94b, #dcb93b)',
        yellowPanelColor: '#505050',
//        panel1bg: 'linear-gradient(to right bottom, #f8deda, #b89e9a)',
        panel2bg: '#fafafa', //'linear-gradient(to right bottom, #f1f1f1, #c1c1c1)',
        panel2color: '#505050',
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