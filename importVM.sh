# Import VMware PowerCLI module
Import-Module VMware.PowerCLI

# Disable SSL certificate warnings (if necessary)
Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Confirm:$false

# OVF file path and VM settings
$OVFPath = "C:\path\to\your-ovf-file.ovf"    # Đường dẫn file OVF
$VMName = "DeployedVM"                      # Tên máy ảo
$Datastore = "datastore1"                   # Datastore để lưu máy ảo
$Network = "VM Network"                     # Tên mạng ảo cần kết nối

# List of ESXi hosts to deploy to
$ESXiHosts = @("esxi-host1.domain", "esxi-host2.domain", "esxi-host3.domain")
$Username = "root"                          # Tên đăng nhập ESXi
$Password = "your-password"                 # Mật khẩu ESXi

# Loop to deploy the OVF to each ESXi host
foreach ($Host in $ESXiHosts) {
    Write-Host "Connecting to ESXi host: $Host..." -ForegroundColor Cyan

    # Connect to ESXi host
    $Connection = Connect-VIServer -Server $Host -User $Username -Password $Password -ErrorAction Stop

    # Get the VMHost object (the ESXi server)
    $VMHost = Get-VMHost -Server $Connection

    # Deploy the OVF template to the current ESXi host
    Write-Host "Deploying OVF to $Host..." -ForegroundColor Green
    Import-VApp -Source $OVFPath `
                -VMHost $VMHost `
                -Name $VMName `
                -Datastore $Datastore `
                -DiskStorageFormat Thin `
                -NetworkMapping @{"Network 1" = $Network}

    Write-Host "Deployment to $Host completed successfully!" -ForegroundColor Yellow

    # Disconnect from the ESXi host
    Disconnect-VIServer -Server $Host -Confirm:$false
}
