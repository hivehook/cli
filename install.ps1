<#
  HiveHook CLI installer for Windows.
    irm https://hivehook.com/install.ps1 | iex

  Environment overrides:
    HIVEHOOK_INSTALL   install directory (default: %LOCALAPPDATA%\HiveHook\bin)
    HIVEHOOK_VERSION   version to install (default: latest release)
#>
$ErrorActionPreference = 'Stop'
$Repo = 'hivehook/cli'
$Bin  = 'hivehook'

$e = [char]27
if ($Host.UI.SupportsVirtualTerminal -ne $false -and -not $env:NO_COLOR) {
	$Bold="$e[1m"; $Reset="$e[0m"; $Rose="$e[38;2;244;63;94m"
	$Green="$e[38;5;42m"; $Grey="$e[38;5;245m"
} else {
	$Bold=''; $Reset=''; $Rose=''; $Green=''; $Grey=''
}
$Check = "$Green$([char]0x2713)$Reset"

function Step($msg) { Write-Host "   $Check $msg" }
function Fail($msg) { Write-Host "   $e[38;5;203m$([char]0x2717) $msg$Reset"; exit 1 }

function Banner {
	Write-Host ""
	Write-Host "   $Rose$([char]0x2B22)$Reset  ${Bold}HiveHook$Reset"
	Write-Host "      ${Grey}Webhook infrastructure CLI$Reset"
	Write-Host ""
}

function Render-Bar($now, $total, $w) {
	$pct = if ($total -gt 0) { [math]::Min(100, [math]::Floor($now * 100 / $total)) } else { 0 }
	$fill = [math]::Floor($pct * $w / 100)
	$filled = [string]([char]0x2B22) * $fill
	$empty  = [string]([char]0x2B21) * ($w - $fill)
	$line = "`r   ${Bold}Downloading$Reset $Rose$filled$Grey$empty$Reset $Bold{0,3}%$Reset" -f $pct
	Write-Host -NoNewline $line
}

function Download-WithBar($url, $dest) {
	$req = [System.Net.HttpWebRequest]::Create($url)
	$req.AllowAutoRedirect = $true
	$req.UserAgent = 'hivehook-installer'
	$resp = $req.GetResponse()
	$total = $resp.ContentLength
	$in = $resp.GetResponseStream()
	$out = [System.IO.File]::Create($dest)
	try {
		$buf = New-Object byte[] 131072
		$readTotal = 0
		while (($n = $in.Read($buf, 0, $buf.Length)) -gt 0) {
			$out.Write($buf, 0, $n)
			$readTotal += $n
			Render-Bar $readTotal $total 24
		}
		Render-Bar $readTotal $total 24
	} finally {
		$out.Close(); $in.Close(); $resp.Close()
	}
	Write-Host ""
}

function Resolve-Version {
	if ($env:HIVEHOOK_VERSION) { return $env:HIVEHOOK_VERSION }
	$rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" `
		-Headers @{ 'User-Agent' = 'hivehook-installer' }
	if (-not $rel.tag_name) { Fail 'could not resolve latest release' }
	return $rel.tag_name
}

function Main {
	Banner
	[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12

	$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
		'AMD64' { 'amd64' }
		'ARM64' { 'arm64' }
		default { Fail "unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)" }
	}
	$platform = "windows_$arch"
	Step "Detected ${Bold}$platform$Reset"

	$version = Resolve-Version
	Step "Installing ${Bold}$Bin $version$Reset"

	$v = $version.TrimStart('v')
	$asset = "${Bin}_${v}_${platform}.zip"
	$base = "https://github.com/$Repo/releases/download/$version"
	$tmp = Join-Path $env:TEMP ("hivehook_" + [guid]::NewGuid().ToString('N'))
	New-Item -ItemType Directory -Path $tmp | Out-Null
	try {
		Download-WithBar "$base/$asset" "$tmp\$asset"

		try {
			Invoke-WebRequest -Uri "$base/checksums.txt" -OutFile "$tmp\checksums.txt" -UseBasicParsing
			$want = (Select-String -Path "$tmp\checksums.txt" -Pattern ([regex]::Escape($asset)) `
				| Select-Object -First 1).Line.Split(' ')[0]
			$got = (Get-FileHash "$tmp\$asset" -Algorithm SHA256).Hash.ToLower()
			if ($want -and ($want -ne $got)) { Fail 'checksum mismatch' }
			Step 'Verified checksum'
		} catch { }

		Expand-Archive -Path "$tmp\$asset" -DestinationPath $tmp -Force

		$dest = if ($env:HIVEHOOK_INSTALL) { $env:HIVEHOOK_INSTALL } else { "$env:LOCALAPPDATA\HiveHook\bin" }
		New-Item -ItemType Directory -Path $dest -Force | Out-Null
		Copy-Item "$tmp\$Bin.exe" "$dest\$Bin.exe" -Force
		Step "Installed to ${Bold}$dest\$Bin.exe$Reset"

		$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
		if ($userPath -notlike "*$dest*") {
			[Environment]::SetEnvironmentVariable('Path', "$userPath;$dest", 'User')
			Write-Host "   ${Grey}Added $dest to your PATH (restart your terminal).$Reset"
		}
	} finally {
		Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
	}

	Write-Host ""
	Write-Host "   $Green${Bold}HiveHook CLI is ready.$Reset"
	Write-Host "   Run ${Bold}hivehook login$Reset to get started."
	Write-Host ""
}

Main
