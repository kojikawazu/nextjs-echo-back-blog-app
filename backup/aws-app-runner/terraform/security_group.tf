# # ---------------------------------------------
# # Security Group
# # ---------------------------------------------
# # App Runner Security Group
# resource "aws_security_group" "app_runner_sg" {
#   name        = "${var.project}-${var.environment}-app-run-sg"
#   description = "app runner security group"
#   vpc_id      = aws_vpc.vpc.id

#   tags = {
#     Name    = "${var.project}-${var.environment}-app-run-sg"
#     Project = var.project
#     Env     = var.environment
#   }
# }

# resource "aws_security_group_rule" "app_runner_in_8080" {
#   security_group_id = aws_security_group.app_runner_sg.id
#   type              = "ingress"
#   protocol          = "tcp"
#   from_port         = var.api_port
#   to_port           = var.api_port
#   cidr_blocks       = [var.default_route]
# }

# resource "aws_security_group_rule" "app_runner_in_db" {
#   security_group_id = aws_security_group.app_runner_sg.id
#   type              = "ingress"
#   protocol          = "tcp"
#   from_port         = 6543
#   to_port           = 6543
#   cidr_blocks       = [var.default_route]
# }

# resource "aws_security_group_rule" "app_runner_out_internet" {
#   security_group_id = aws_security_group.app_runner_sg.id
#   type              = "egress"
#   protocol          = "-1"
#   from_port         = 0
#   to_port           = 0
#   cidr_blocks       = [var.default_route]
# }
