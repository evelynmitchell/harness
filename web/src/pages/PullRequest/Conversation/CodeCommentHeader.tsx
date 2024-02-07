import { Button, ButtonVariation, Container, Layout, Text, Utils } from '@harnessio/uicore'
interface CodeCommentHeaderProps extends Pick<GitInfoProps, 'repoMetadata' | 'pullReqMetadata'> {
  pullReqMetadata
          <Layout.Horizontal flex={{ alignItems: 'center' }}>
            <Text
              inline
              className={css.fname}
              lineClamp={1}
              tooltipProps={{
                portalClassName: css.popover
              }}>
              <Link
                // className={css.fname}
                to={`${routes.toCODEPullRequest({
                  repoPath: repoMetadata?.path as string,
                  pullRequestId: String(pullReqMetadata?.number),
                  pullRequestSection: PullRequestSection.FILES_CHANGED
                })}?path=${commentItems[0].payload?.code_comment?.path}&commentId=${commentItems[0].payload?.id}`}>
                {commentItems[0].payload?.code_comment?.path}
              </Link>
            </Text>
            <Button
              variation={ButtonVariation.ICON}
              icon="code-copy"
              className={css.copyButton}
              iconProps={{ size: 14 }}
              onClick={() => {
                if (commentItems[0].payload?.code_comment?.path) {
                  Utils.copy(commentItems[0].payload?.code_comment?.path)
                }
              }}
            />
          </Layout.Horizontal>