
Write-Host "----------- Test GET request -----------" 


$Url = "http://localhost:8080/passport?series=8003&number=011384"

# simple html
Write-Host 'request as ContentType "text/html"' 

$response = Invoke-RestMethod -Method GET -Uri $url -ContentType "text/html"

Write-Host "Response:" 
Write-Host $response

# JSON response in body
Write-Host 'request as ContentType "application/json"' 

$response = Invoke-RestMethod -Method GET -Uri $url -ContentType "application/json"

Write-Host "JSON in Response:" 

Write-Host ($response | ConvertTo-Json ) 



# 
Write-Host "----------- Test POST request -----------" 

$Url = "http://localhost:8080/passport"

#массив для хранения списка паспортов 
$passport_list = New-Object System.Collections.ArrayList
# nonvalid
$passport = @{
	series = "8003"
	number = "011384"
}
$cnt =  $passport_list.add($passport)

# valid
$passport = @{
	series = "4050"
	number = "039589"
}
$cnt = $passport_list.add($passport)

# nonvalid
$passport = @{
	series = "5203"
	number = "257719"
}
$cnt = $passport_list.add($passport)
# nonvalid
$passport = @{
	series = "5000"
	number = "347024"
}
$cnt = $passport_list.add($passport)

# nonvalid
$passport = @{
	series = "2507"
	number = "857721"
}
$cnt = $passport_list.add($passport)

# nonvalid
$passport = @{
	series = "2507"
	number = "857728"
}
$cnt = $passport_list.add($passport)



 
$json_body =  $passport_list | ConvertTo-Json

 $params = @{
    Uri         = $Url
    Method      = 'POST'
    Body        = $json_body
    ContentType = 'application/json'
  }

Write-Host "JSON in Request:" 
Write-Host  $json_body



$response = try { Invoke-RestMethod @params }  

catch [Exception] {
            Write-Host "StatusCode:" $_.Exception.Response.StatusCode.value__ 
            Write-Host "StatusDescription:" $_.Exception.Response.StatusDescription
        }

Write-Host "JSON in Response:" 

$response = ($response | ConvertTo-Json)
Write-Host $response

#foreach ($passport in $response) {
#	Write-Host $passport
#}

# 
#Start-Sleep -s 5


