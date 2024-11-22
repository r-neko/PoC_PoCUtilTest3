package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func initLocalRepo(repoName, license, orgName string) error {
	// カレントディレクトリをリポジトリ化
	repo, err := git.PlainInit(".", false) // カレントディレクトリを指定
	if err != nil {
		return fmt.Errorf("failed to initialize repository in current directory: %w", err)
	}
	fmt.Println("Initialized empty Git repository in the current directory")

	// ファイル作成（READMEとLICENSE）
	err = os.WriteFile("README.md", []byte("# "+repoName), 0644)
	if err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}
	licenseContent, err := getLicenseContent(license)
	if err != nil {
		return fmt.Errorf("failed to get license content: %w", err)
	}
	err = os.WriteFile("LICENSE", []byte(licenseContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create LICENSE: %w", err)
	}

	// 作業ツリーの取得
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// ファイルをステージング
	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// 初期コミット
	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your.email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// リモート設定
	remoteURL := fmt.Sprintf("https://github.com/%s/%s.git", orgName, repoName)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	// トークンを環境変数から取得
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// リモートへPush
	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "username", // GitHubでは任意の値
			Password: token,      // トークンをパスワードとして利用
		},
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Println("Local repository initialized and pushed successfully.")
	return nil
}
