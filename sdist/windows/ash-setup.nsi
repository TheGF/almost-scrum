;NSIS Modern User Interface
;Start Menu Folder Selection Example Script
;Written by Joost Verburg

;--------------------------------
;Include Modern UI

  !include "MUI2.nsh"

;--------------------------------
;General

  ;Name and file
  Name "Almost Scrum"
  OutFile "..\..\dist\almost-scrum-setup-0.5.exe"
  Unicode True

  ;Default installation folder
  InstallDir "$LOCALAPPDATA\AlmostScrum"
  
  ;Get installation folder from registry if available
  InstallDirRegKey HKCU "Software\AlmostScrum" ""

  ;Request application privileges for Windows Vista
  RequestExecutionLevel user

  ShowInstDetails Show

;--------------------------------
;Variables

  Var StartMenuFolder

;--------------------------------
;Interface Settings

  !define MUI_ABORTWARNING

;--------------------------------
;Pages

  !insertmacro MUI_PAGE_LICENSE "License.txt"
  !insertmacro MUI_PAGE_COMPONENTS
  !insertmacro MUI_PAGE_DIRECTORY
  
  ;Start Menu Folder Page Configuration
  !define MUI_STARTMENUPAGE_REGISTRY_ROOT "HKCU" 
  !define MUI_STARTMENUPAGE_REGISTRY_KEY "Software\AlmostScrum" 
  !define MUI_STARTMENUPAGE_REGISTRY_VALUENAME "Start Menu Folder"
  
  !insertmacro MUI_PAGE_STARTMENU Application $StartMenuFolder
  
  !insertmacro MUI_PAGE_INSTFILES
  
  !insertmacro MUI_UNPAGE_CONFIRM
  !insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
;Languages
 
  !insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Installer Sections

Section "Almost Scrum" AlmostScrum

  SetOutPath "$INSTDIR"
  
  ;ADD YOUR OWN FILES HERE...
  File /oname=ash.exe ..\bin\ash_windows.exe 
  File License.txt 

  ;Store installation folder
  WriteRegStr HKCU "Software\AlmostScrum" "" $INSTDIR
  
  ;Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; Set back to HKCU
  EnVar::SetHKCU
  DetailPrint "EnVar::SetHKCU"
 
    ; Add an expanded value
  EnVar::AddValue "PATH" "$INSTDIR"
  Pop $0
  DetailPrint "EnVar::AddValue returned=|$0|"

;   EnVar::AddValueEx "path" "$INSTDIR"
;   Pop $0
;   DetailPrint "EnVar::AddValue returned=|$0|"

  !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
    
    CreateDirectory "$DOCUMENTS\AlmostScrum"
    ;Create shortcuts
    CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
    CreateShortCut "$SMPROGRAMS\$StartMenuFolder\Almost Scrum.lnk" "$INSTDIR\ash.exe" \
        "server $DOCUMENTS\AlmostScrum" "$INSTDIR\ash.exe" 1 SW_SHOWMINIMIZED \
        ALT|CONTROL|F8 "Almost Scrum Server"
    CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk" "$INSTDIR\Uninstall.exe" 
  
  !insertmacro MUI_STARTMENU_WRITE_END

SectionEnd

;Uninstaller Section

Section "Uninstall"

  Delete "$INSTDIR\ash.exe"

  Delete "$INSTDIR\Uninstall.exe"

  RMDir "$INSTDIR"
  
  !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
    
  Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk"
  RMDir "$SMPROGRAMS\$StartMenuFolder"
  
  DeleteRegKey /ifempty HKCU "Software\Almost Scrum"

SectionEnd
    