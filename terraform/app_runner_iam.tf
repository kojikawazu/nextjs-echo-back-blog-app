# ---------------------------------------------
# App Runner - IAMロール
# ---------------------------------------------
resource "aws_iam_role" "apprunner_access_role" {
  name = "${var.project}-${var.environment}-apprunn-access-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : ["build.apprunner.amazonaws.com", "tasks.apprunner.amazonaws.com"]
        },
        "Action" : "sts:AssumeRole"
      }
    ]
  })
}

# ---------------------------------------------
# App Runner - IAMロールポリシー
# ---------------------------------------------
resource "aws_iam_role_policy" "apprunner_access_policy" {
  role = aws_iam_role.apprunner_access_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:CreateLogGroup",
          "logs:DescribeLogStreams",
          "logs:DescribeLogGroups",
          "sts:AssumeRole",
          "ecr:GetAuthorizationToken",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "apprunner:DescribeService",
          "apprunner:ListServices",
          "apprunner:ListTagsForResource",
          "apprunner:TagResource",
          "apprunner:UntagResource",
          "apprunner:UpdateService",
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = "*"
      }
    ]
  })
}
