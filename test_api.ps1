# Test Script for SinkedIn API
# Must have curl installed

# Base URL
$baseUrl = "http://localhost:8080/api"
$token = ""

Write-Host "Testing User Registration and Login..."
# 1. Register a new user
$registerData = @{
    name = "John Doe"
    username = "johndoe"
    email = "john@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/users/register" -Method Post -Body $registerData -ContentType "application/json"
Write-Host "Registration Response:" $response

# 2. Login
$loginData = @{
    email = "john@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/users/login" -Method Post -Body $loginData -ContentType "application/json"
$token = $response.token
Write-Host "Login successful. Token:" $token

# Set default headers for authenticated requests
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# 3. Create a post
Write-Host "`nTesting Post Creation..."
$postData = @{
    content = "Hello world! #firstpost"
    hashtags = @("firstpost")
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/posts" -Method Post -Headers $headers -Body $postData
Write-Host "Created Post:" $response
$postId = $response.id

# 4. Get posts
Write-Host "`nTesting Get Posts..."
$response = Invoke-RestMethod -Uri "$baseUrl/posts" -Method Get -Headers $headers
Write-Host "Posts:" $response

# 5. Create a comment
Write-Host "`nTesting Comment Creation..."
$commentData = @{
    postId = $postId
    content = "Great first post! #welcome"
    type = "normal"
    hashtags = @("welcome")
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/comments" -Method Post -Headers $headers -Body $commentData
Write-Host "Created Comment:" $response
$commentId = $response.id

# 6. Like a post
Write-Host "`nTesting Like Functionality..."
$response = Invoke-RestMethod -Uri "$baseUrl/likes/post/$postId" -Method Post -Headers $headers
Write-Host "Like Response:" $response

# 7. Get post likes
$response = Invoke-RestMethod -Uri "$baseUrl/likes/post/$postId" -Method Get -Headers $headers
Write-Host "Post Likes:" $response

# 8. Get trending hashtags
Write-Host "`nTesting Trending Hashtags..."
$response = Invoke-RestMethod -Uri "$baseUrl/hashtags/trending" -Method Get -Headers $headers
Write-Host "Trending Hashtags:" $response

# 9. Register another user for follow testing
$register2Data = @{
    name = "Jane Smith"
    username = "janesmith"
    email = "jane@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/users/register" -Method Post -Body $register2Data -ContentType "application/json"
Write-Host "`nSecond User Registration:" $response

# 10. Follow user
Write-Host "`nTesting Follow Functionality..."
$response = Invoke-RestMethod -Uri "$baseUrl/follow/janesmith" -Method Post -Headers $headers
Write-Host "Follow Response:" $response

# 11. Get followers
$response = Invoke-RestMethod -Uri "$baseUrl/follow/followers/janesmith" -Method Get -Headers $headers
Write-Host "Followers:" $response

Write-Host "`nTest Script Completed!"
